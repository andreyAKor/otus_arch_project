<?php
declare(strict_types=1);

namespace common\components\services;

use common\components\HttpStatus;
use yii\base\BaseObject;
use yii\web\ServerErrorHttpException;

/**
 * Класс-компонент для работы с API сервисов
 * @package common\components\services
 */
final class Service extends BaseObject
{
    /**
     * Тайм-аут ожидания запроса, в секундах
     * @var integer
     */
    public $requestTimeout = 20;

    /**
     * URL адрес gateway-сервисов
     * @var string
     */
    public string $url = 'http://0.0.0.0:6080';

    /**
     * Возвращает указатель на клиента
     * @return Client
     */
    private function getClient(): Client
    {
        return new Client([
            'requestTimeout' => $this->requestTimeout,
            'url' => $this->url
        ]);
    }

    /**
     * Создание заявки
     * @param int $coinType - тип монет
     * @param string $value - отдаваемое значение
     * @param string $address - адрес получения обмена
     * @return array|null
     * @throws ServerErrorHttpException
     * @throws \yii\httpclient\Exception
     */
    public function create(int $coinType, string $value, string $address): ?array
    {
        try {
            $response = $this->getClient()->create($coinType, $value, $address);
        } catch (\Throwable $e) {
            throw new ServerErrorHttpException($e->getMessage());
        }

        if (in_array($response->getStatusCode(), [HttpStatus::CREATED, HttpStatus::BAD_REQUEST])) {
            return $response->getData();
        }

        return null;
    }

    /**
     * Возвращает данные заявки
     * @param string $bidId - ID-заявки
     * @return array|null
     * @throws ServerErrorHttpException
     * @throws \yii\httpclient\Exception
     */
    public function get(string $bidId): ?array
    {
        try {
            $response = $this->getClient()->get($bidId);
        } catch (\Throwable $e) {
            throw new ServerErrorHttpException($e->getMessage());
        }

        if ($response->getStatusCode() == HttpStatus::OK) {
            return $response->getData();
        }

        return null;
    }
}
