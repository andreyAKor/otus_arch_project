<?php
declare(strict_types=1);

namespace common\components\monolog;

use Mero\Monolog\MonologComponent as MonologMonologComponent;
use Yii;

/**
 * Класс-компонент для обвёртки класса-компонента MonologComponent
 * @package common\components\monolog
 */
final class MonologComponent extends MonologMonologComponent
{
    /**
     * @var array
     */
    public $channels;

    /**
     * Инициализация компонента
     * @throws \Mero\Monolog\Exception\InvalidHandlerException
     * @throws \yii\base\InvalidConfigException
     */
    public function init(): void
    {
        foreach ($this->channels as $name => &$channel) {
            if (!empty($channel['handler']) && is_array($channel['handler'])) {
                foreach ($channel['handler'] as &$handlerConfig) {
                    if (is_array($handlerConfig) && isset($handlerConfig['class'])) {
                        $handlerConfig = Yii::createObject($handlerConfig);
                    }
                }
            }
        }

        parent::init();
    }
}
