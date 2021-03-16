<?php
declare(strict_types=1);

use yii\db\Migration;

/**
 * Class m200417_111020_CreateServicesLog
 */
class m200417_111020_CreateServicesLog extends Migration
{
    /**
     * @return bool
     * @throws \yii\db\Exception
     */
    public function up(): bool
    {
        $this->db->createCommand('
CREATE TABLE `ServicesLog` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT \'Первичный ключ\',
  `level` int(11) NOT NULL COMMENT \'Уровень записи\',
  `sessionId` varchar(50) NOT NULL COMMENT \'ID сессии\',
  `message` text DEFAULT NULL COMMENT \'Сообщение\',
  `url` varchar(255) NOT NULL COMMENT \'Запрос\',
  `data` json DEFAULT NULL COMMENT \'Параметры запроса\',
  `content` json DEFAULT NULL COMMENT \'Данные запроса\',
  `request` longtext NOT NULL COMMENT \'HTTP запрос\',
  `response` longtext DEFAULT NULL COMMENT \'HTTP ответ\',
  `createdAt` datetime NOT NULL COMMENT \'Создано\',
  PRIMARY KEY (`id`),
  KEY `ServicesLog_level_IDX` (`level`),
  KEY `ServicesLog_url_IDX` (`url`),
  KEY `ServicesLog_createdAt_IDX` (`createdAt`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT=\'Лог обращений к API сервисами\'
')->execute();

        return true;
    }

    /**
     * @return bool
     */
    public function down(): bool
    {
        echo "m200417_111020_CreateServicesLog cannot be reverted.\n";

        return false;
    }
}
