<?php
declare(strict_types = 1);

use yii\db\Migration;

/**
 * Class m200424_051731_UpdServicesLog
 */
class m200424_051731_UpdServicesLog extends Migration
{
    /**
     * @return bool
     * @throws \yii\db\Exception
     */
    public function up(): bool
    {
        $this->db->createCommand('
RENAME TABLE ServicesLog TO services_api_Log;
')->execute();

        return true;
    }

    /**
     * @return bool
     */
    public function down(): bool
    {
        echo "m200424_051731_UpdServicesLog cannot be reverted.\n";

        return false;
    }
}
