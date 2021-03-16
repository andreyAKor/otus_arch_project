<?php
declare(strict_types=1);

namespace frontend\controllers;

use common\helpers\coin\Convert;
use frontend\models\BidForm;
use Yii;
use yii\web\Controller;
use yii\web\ErrorAction;

/**
 * Class SiteController
 * @package frontend\controllers
 */
class SiteController extends Controller
{
    /**
     * {@inheritdoc}
     * @see \yii\base\Controller::actions()
     */
    public function actions(): array
    {
        return [
            'error' => [
                'class' => ErrorAction::class
            ]
        ];
    }

    /**
     * Главная страница
     * @return string
     */
    public function actionIndex(): string
    {
        $model = new BidForm();
        if ($model->load(Yii::$app->request->post())) {
            if ($model->validate()) {
                if ($res = $this->createBid($model)) {
                    $this->redirect([
                        'view',
                        'bidId' => $res['data']['id']
                    ]);
                }
            } else {
                Yii::$app->session->setFlash('error', $model->getFirstErrors()[0]);
            }
        }

        return $this->render('index', [
            'model' => $model
        ]);
    }

    /**
     * Детали завки
     * @param string $bidId
     * @return string
     */
    public function actionView(string $bidId): string
    {
        if (!($data = $this->getBid($bidId))) {
            Yii::$app->end(0, $this->redirect([
                'site/index'
            ]));
        }

        return $this->render('view', [
            'bidId' => $bidId,
            'data' => $data
        ]);
    }

    /**
     * @param BidForm $model
     * @return array|null
     */
    private function createBid(BidForm $model): ?array
    {
        $value = Convert::denormalize($model->coinType, $model->value);

        try {
            // Создание заявки
            if (!($res = Yii::$app->services->create($model->coinType, $value, $model->address))) {
                Yii::$app->session->setFlash('error', 'Internal error (code 1)');
                return null;
            }

            if (isset($res['error'])) {
                Yii::$app->session->setFlash('error', $res['error']);
                return null;
            }

            return $res;
        } catch (\Throwable $e) {
            Yii::$app->session->setFlash('error', 'Internal error (code 2)');
        }

        return null;
    }

    /**
     * @param string $bidId
     * @return array|null
     */
    private function getBid(string $bidId): ?array
    {
        try {
            // Возвращает данные заявки
            if (!($res = Yii::$app->services->get($bidId))) {
                Yii::$app->session->setFlash('error', 'Internal error (code 3)');
                return null;
            }

            return $res;
        } catch (\Throwable $e) {
            Yii::$app->session->setFlash('error', 'Internal error (code 4)');
        }

        return null;
    }
}
