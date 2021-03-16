<?php
declare(strict_types=1);

namespace common\components\monolog\handler;

/**
 * Класс для логирования общения с API через обвёртку yii2-monolog
 * @package common\components\monolog\handler
 */
final class Api extends AbstractHandler
{
    /**
     * Пишем лог
     * @param array $record - данные для записи в лог
     * @throws \yii\db\Exception
     */
    protected function write(array $record): void
    {
        $db = $this->getYiiConnection($this->reference);

        $db->quoteTableName($this->table);
        $db->createCommand()
            ->insert($this->table, array_merge([
                'level' => $record['level'],
                'message' => $record['message'] ? $record['message'] : null,
                'createdAt' => $record['datetime']->format('Y-m-d H:i:s')
            ], $record['context']))
            ->execute();
    }
}
