<?php
declare(strict_types=1);

use common\helpers\coin\CoinType;
use yii\bootstrap\ActiveForm;
use yii\helpers\Html;

/* @var $model \frontend\models\BidForm */

$this->title = 'Создание заявки';
$this->params['breadcrumbs'][] = $this->title;
?>
<div class="site-index">
    <h1><?= Html::encode($this->title) ?></h1>
    <div class="mar-bot-13 clearfix">
        <div class="mar-bot-13 clearfix">
            <div class="col-sm-2">
                <?php
                $form = ActiveForm::begin([
                    'method' => 'post',
                ]);
                ?>
            </div>
        </div>
        <div class="mar-bot-13 clearfix">
            <div class="col-sm-2"><?= $form->field($model, 'coinType')->dropDownList(CoinType::getList()); ?></div>
        </div>
        <div class="mar-bot-13 clearfix">
            <div class="col-sm-2"><?= $form->field($model, 'value')->textInput(['maxlength' => 64]); ?></div>
        </div>
        <div class="mar-bot-13 clearfix">
            <div class="col-sm-2"><?= $form->field($model, 'address')->textInput(['maxlength' => 128]); ?></div>
        </div>
        <div class="mar-bot-13 clearfix">
            <div class="col-sm-2"><?= Html::submitButton('Создать', ['class' => 'btn btn-success']) ?></div>
        </div>
        <?php ActiveForm::end(); ?>
    </div>
</div>
