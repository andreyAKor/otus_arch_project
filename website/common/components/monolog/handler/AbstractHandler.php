<?php
declare(strict_types=1);

namespace common\components\monolog\handler;

use Carbon\Carbon;
use Monolog\Handler\AbstractProcessingHandler;
use Yii;
use yii\db\Connection;
use yii\di\Instance;

/**
 * Абстрактный класс для логирования общения с через обвёртку yii2-monolog
 * @package common\components\monolog\handler
 */
abstract class AbstractHandler extends AbstractProcessingHandler
{
    /**
     * Время жизни логов в часах
     * @var int
     */
    public int $logLiferime = PHP_INT_MAX;

    /**
     * Yii2 database connection
     * @var string
     */
    public string $reference = 'db';

    /**
     * Log table in database
     * @var string
     */
    public string $table = '';

    /**
     * Деструктор класса
     */
    public function __destruct()
    {
        $db = $this->getYiiConnection($this->reference);

        $db->quoteTableName($this->table);
        $db->createCommand()
            ->delete($this->table, 'createdAt <= :createdAt', [
                ':createdAt' => Carbon::now(Yii::$app->timeZone)->subHours($this->logLiferime)
                    ->format('Y-m-d H:i:s')
            ])
            ->execute();

        parent::__destruct();
    }

    /**
     * Return a Yii2 database connection.
     * @param string $reference Name of Yii2 database connection
     * @return Connection
     * @throws \yii\base\InvalidConfigException
     */
    protected function getYiiConnection($reference): Connection
    {
        return Instance::ensure($reference, Connection::className());
    }
}
