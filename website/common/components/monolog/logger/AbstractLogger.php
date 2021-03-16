<?php
declare(strict_types=1);

namespace common\components\monolog\logger;

/**
 * Абстрактный класс для описания формирования логирования
 * @package common\components\monolog\logger
 */
abstract class AbstractLogger
{
    /**
     * Возвращает данные для логирования
     * @return array
     */
    public function context(): array
    {
        return $this->_context;
    }
}
