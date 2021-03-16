<?php
declare(strict_types=1);

use yii\db\Migration;

/**
 * Class m191210_180026_UpdDataOfUserAccount
 */
class m191210_180026_UpdDataOfUserAccount extends Migration
{
    /**
     * @return bool
     * @throws \yii\db\Exception
     */
    public function up(): bool
    {
        $this->db->createCommand('
INSERT INTO `user` (id, username, auth_key, password_hash, password_reset_token, email, status, created_at, updated_at, verification_token) VALUES(2, \'admin\', \'D2MjIQADPyFrMf2TaFJq9AtTwsWF83q2\', \'$2y$13$o3VW/WqICYyF..Y3IKiBX.tXLGJx1i0ti5ezsAZrrvTOv3tasMafu\', NULL, \'info@webumka.ru\', 10, 1575472876, 1575472918, \'NS98q5lF1IJk-L3TxsIPBU8Ta8Wlvimk_1575472876\');
')->execute();

        return true;
    }

    /**
     * @return bool
     */
    public function down(): bool
    {
        echo "m191210_180026_UpdDataOfUserAccount cannot be reverted.\n";

        return false;
    }
}
