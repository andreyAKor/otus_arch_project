<?php
declare(strict_types=1);

namespace frontend\models;

use common\helpers\coin\CoinType;
use yii\base\Model;

/**
 * Class BidForm
 * @package frontend\models
 */
class BidForm extends Model
{
    /**
     * @var int
     */
    public int $coinType = CoinType::ETH;

    /**
     * @var string
     */
    public string $value = '';

    /**
     * @var string
     */
    public string $address = '';

    /**
     * @return array[]
     */
    public function rules(): array
    {
        return [
            [
                [
                    'coinType',
                    'value',
                    'address'
                ],
                'required'
            ],

            // value
            [
                'value',
                'required'
            ],
            [
                'value',
                'string',
                'max' => 64
            ],

            // address
            [
                'address',
                'required'
            ],
            [
                'address',
                'string',
                'max' => 128
            ],
        ];
    }

    /**
     * @return array
     */
    public function attributeLabels(): array
    {
        return [
            'coinType' => 'Какой тип монет вы отдаете?',
            'value' => 'Какое количество монет хотите отдать?',
            'address' => 'На какой адрес Вам прислать обменянные монеты?'
        ];
    }
}
