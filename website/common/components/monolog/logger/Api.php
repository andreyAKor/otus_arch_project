<?php
declare(strict_types=1);

namespace common\components\monolog\logger;

use Yii;
use yii\httpclient\Request;
use yii\httpclient\Response;

/**
 * Класс-хелпер для подготовки логгирования данных общения с API
 * @package common\components\monolog\logger
 */
final class Api extends AbstractLogger
{
    /**
     * @var array
     */
    protected array $_context = [];

    /**
     * @param Request $request - объект Request запроса
     * @param Response|null $response - объект Response ответа
     * @return Api
     */
    public static function factory(Request $request, ?Response $response): Api
    {
        $self = new self();
        $self->_context = [
            'sessionId' => Yii::$app->session->getId(),
            'url' => $request->getUrl(),
            'data' => $request->getData(),
            'content' => print_r($request->getContent(), true),

            'request' => (string)$request,
            'response' => (string)$response
        ];

        return $self;
    }
}
