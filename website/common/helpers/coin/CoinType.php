<?php
declare(strict_types=1);

namespace common\helpers\coin;

/**
 * Class CoinType
 * @package common\helpers\coin
 */
final class CoinType
{
    /**
     * Эфир
     * @var integer
     */
    public const ETH = 1;

    /**
     * Битки
     * @var integer
     */
    public const BTC = 2;

    /**
     * Сопоставление типов монет с их именам
     *
     * @var array
     */
    private static $_coinTypeNames = [
        self::ETH => 'ETH',
        self::BTC => 'BTC'
    ];

    /**
     * Возвращает имя монеты зная ее тип
     *
     * @param int $coinType
     * @return string|NULL
     */
    public static function getName(int $coinType): ?string
    {
        return isset(self::$_coinTypeNames[$coinType]) ? self::$_coinTypeNames[$coinType] : null;
    }

    /**
     * Возвращает список типов валют и их имена
     *
     * @return array
     */
    public static function getList(): array
    {
        return [
            self::ETH => self::getName(self::ETH),
            self::BTC => self::getName(self::BTC)
        ];
    }
}