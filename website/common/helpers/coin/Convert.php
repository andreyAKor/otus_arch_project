<?php
declare(strict_types=1);

namespace common\helpers\coin;

/**
 * Class Convert
 * @package common\helpers\coin
 */
class Convert
{
    /**
     * Возвращает значение монеты в нормализованном формате
     * @param int $coinType
     * @param string $value
     * @param int|null $len
     * @return string|null
     */
    public static function normalize(int $coinType, string $value, int $len = null): ?string
    {
        if (!($divisor = self::getDivisor($coinType))) {
            return null;
        }

        $res = bcdiv($value, $divisor, strlen($divisor));
        $res = rtrim($res, '0');

        if ($len !== null && ($pos = strpos($res, '.'))) {
            $res = substr($res, 0, $pos) . '.' . substr($res, $pos + 1, $len);
            $res = rtrim($res, '0');
        }

        if (substr($res, -1) == '.') {
            $res = substr($res, 0, strlen($res) - 1);
        }

        return $res;
    }

    /**
     * Возвращает значение монеты в денормализованном формате
     * @param int $coinType
     * @param string $value
     * @return string|null
     */
    public static function denormalize(int $coinType, string $value): ?string
    {
        if (!($divisor = self::getDivisor($coinType))) {
            return null;
        }

        $res = bcmul($value, $divisor, strlen($divisor));
        $res = rtrim($res, '0');

        if (substr($res, -1) == '.') {
            $res = substr($res, 0, strlen($res) - 1);
        }

        return $res;
    }

    /**
     * @param int $coinType
     * @return string|null
     */
    public static function getDivisor(int $coinType): ?string
    {
        switch ($coinType) {
            case CoinType::ETH:
                return '1000000000000000000';
            case CoinType::BTC:
                return '100000000';
        }

        return null;
    }
}