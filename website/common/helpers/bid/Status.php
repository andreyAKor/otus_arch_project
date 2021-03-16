<?php
declare(strict_types=1);

namespace common\helpers\bid;

/**
 * Class Status
 * @package common\helpers\order
 */
final class Status
{
    /**
     * Ожидание отправки
     * @var string
     */
    public const PENDING = 1;

    /**
     * Закрыт
     * @var string
     */
    public const CLOSE = 2;

    /**
     * Ошибка
     * @var int
     */
    public const ERROR = 3;

    /**
     * Сопоставление кодов статусов их именам
     * @var array
     */
    private static array $_statusNames = [
        self::PENDING => 'Ожидание отправки',
        self::CLOSE => 'Закрыт',
        self::ERROR => 'Ошибка'
    ];

    /**
     * Возвращает имя статуса по его коду
     * @param string $status
     * @return string|NULL
     */
    public static function getName(string $status): ?string
    {
        return isset(self::$_statusNames[$status]) ? self::$_statusNames[$status] : null;
    }
}