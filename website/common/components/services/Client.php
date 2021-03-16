<?php
declare(strict_types=1);

namespace common\components\services;

use common\components\monolog\logger\Api;
use Yii;
use yii\base\BaseObject;
use yii\base\Exception;
use yii\httpclient\Client as httpClient;
use yii\httpclient\CurlTransport;
use yii\httpclient\Request;
use yii\httpclient\Response;

/**
 * Class Client
 * @package common\components\services
 */
final class Client extends BaseObject
{
    private const USER_AGENT = 'ServicesClient/1.0.0';

    /**
     * Тайм-аут ожидания запроса, в секундах
     * @var integer
     */
    public int $requestTimeout = 20;

    /**
     * URL адрес API сервисов
     * @var string
     */
    public string $url = 'http://0.0.0.0:6080';

    /**
     * Создание заявки
     * @param int $coinType - тип монет
     * @param string $value - отдаваемое значение
     * @param string $address - адрес получения обмена
     * @return Response
     * @throws Exception
     * @throws \Throwable
     */
    public function create(int $coinType, string $value, string $address): Response
    {
        return $this->request('/create', [
            'coin_type' => $coinType,
            'value' => $value,
            'address' => $address,
        ], 'post');
    }

    /**
     * Возвращает данные заявки
     * @param string $bidId - ID-заявки
     * @return Response
     */
    public function get(string $bidId): Response
    {
        return $this->request('/get', [
            'id' => $bidId
        ]);
    }

    /**
     * Совершает запросы к API
     * @param string $resource - Запрос
     * @param array|null $params - параметры запроса
     * @param string $method - тип запроса
     * @return Response
     * @throws Exception
     * @throws \Throwable
     * @throws \yii\base\InvalidConfigException
     * @throws \yii\httpclient\Exception
     */
    protected function request(string $resource, ?array $params = null, string $method = 'get'): Response
    {
        // Клиент
        $client = new httpClient();
        $client->setTransport(CurlTransport::class);

        // Предварительная подготовка к запросу
        $request = $client->createRequest()
            ->setMethod($method)
            ->setUrl($this->url . $resource)
            ->setHeaders([
                'Accept' => '*/*',
                'Accept-Encoding' => 'deflate',
                'Cache-Control' => 'no-cache',
                'Connection' => 'keep-alive',
                'User-Agent' => self::USER_AGENT,
                'Content-Type' => 'application/x-www-form-urlencoded'
            ])
            ->addOptions([
                'timeout' => $this->requestTimeout
            ]);

        // Если имеются параметры запроса, то задаём их
        if (is_array($params) && count($params)) {
            $request->setData($params);
        }

        $response = null;

        try {
            // Совершаем запрос
            if (!($response = $request->send())) {
                throw new Exception('Ошибка запроса');
            }
        } catch (\Throwable $e) {
            $this->addError($e->getMessage(), $request, $response);

            throw $e;
        }

        if (!$response || $response->getContent() == '') {
            $exMsg = 'Ошибка запроса, statusCode: пусто ответ';

            $this->addError($exMsg, $request, $response);

            throw new Exception($exMsg);
        }

        $this->addInfo($request, $response);

        return $response;
    }

    /**
     * @param string $comment
     * @param Request $request
     * @param Response|null $response
     */
    private function addError(string $comment, Request $request, ?Response $response): void
    {
        $ctx = Api::factory($request, $response)->context();
        Yii::$app->monolog->getLogger('services')->addError($comment, $ctx);
    }

    /**
     * @param Request $request
     * @param Response $response
     */
    private function addInfo(Request $request, Response $response): void
    {
        $ctx = Api::factory($request, $response)->context();
        Yii::$app->monolog->getLogger('services')->addInfo(null, $ctx);
    }
}
