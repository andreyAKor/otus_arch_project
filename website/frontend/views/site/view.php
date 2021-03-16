<?php
declare(strict_types=1);

use common\helpers\bid\Status;
use common\helpers\coin\CoinType;
use common\helpers\coin\Convert;
use yii\helpers\Html;

/* @var $bidId int */
/* @var $data array */

$this->title = 'Детали заявки';
$this->params['breadcrumbs'][] = $this->title;
?>
<div class="site-index">
    <h1><?= Html::encode($this->title) ?></h1>
    <div class="mar-bot-13 clearfix">
        <div class="mar-bot-13 clearfix">
            <p><b>№ заявки:</b> <?= $bidId ?></p>
            <p><b>Статус:</b> <?= Status::getName((string)$data['data']['status']) ?></p>

            <h3>Вы отдаете</h3>
            <p>
                <b>Монеты:</b> <?= Convert::normalize((int)$data['data']['receivedCoinType'], $data['data']['receivedValue'], 18) ?> <?= CoinType::getName((int)$data['data']['receivedCoinType']) ?>
            </p>
            <p><b>На адрес:</b> <?= $data['data']['receivedAddress'] ?></p>

            <h3>Вы получаете</h3>
            <p>
                <b>Монеты:</b> <?= Convert::normalize((int)$data['data']['givenCoinType'], $data['data']['givenValue'], 18) ?> <?= CoinType::getName((int)$data['data']['givenCoinType']) ?>
            </p>
            <p><b>На адрес:</b> <?= $data['data']['givenAddress'] ?></p>
        </div>
    </div>
</div>
