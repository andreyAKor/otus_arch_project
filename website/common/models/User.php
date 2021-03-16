<?php
declare(strict_types=1);

namespace common\models;

use Yii;
use yii\base\NotSupportedException;
use yii\behaviors\TimestampBehavior;
use yii\db\ActiveRecord;
use yii\web\IdentityInterface;

/**
 * Class User
 * @package common\models
 */
class User extends ActiveRecord implements IdentityInterface
{
    /**
     * @var integer
     */
    public const STATUS_DELETED = 0;

    /**
     * @var integer
     */
    public const STATUS_INACTIVE = 9;

    /**
     * @var integer
     */
    public const STATUS_ACTIVE = 10;

    /**
     * @return string
     */
    public static function tableName(): string
    {
        return '{{%user}}';
    }

    /**
     * @return array
     */
    public function behaviors(): array
    {
        return [
            TimestampBehavior::className()
        ];
    }

    /**
     * @return array
     */
    public function rules(): array
    {
        return [
            // status
            [
                'status',
                'default',
                'value' => self::STATUS_INACTIVE
            ],
            [
                'status',
                'in',
                'range' => [
                    self::STATUS_ACTIVE,
                    self::STATUS_INACTIVE,
                    self::STATUS_DELETED
                ]
            ]
        ];
    }

    /**
     * @param int $id
     * @return User|NULL
     */
    public static function findIdentity($id): ?User
    {
        return static::findOne([
            'id' => $id,
            'status' => self::STATUS_ACTIVE
        ]);
    }

    /**
     * @param string $token
     * @param int|null $type
     * @throws NotSupportedException
     */
    public static function findIdentityByAccessToken($token, $type = null): void
    {
        throw new NotSupportedException('"findIdentityByAccessToken" is not implemented.');
    }

    /**
     * Finds user by username
     * @param string $username
     * @return User|NULL
     */
    public static function findByUsername(string $username): ?User
    {
        return static::findOne([
            'username' => $username,
            'status' => self::STATUS_ACTIVE
        ]);
    }

    /**
     * Finds user by password reset token
     * @param string $token - password reset token
     * @return User|null
     */
    public static function findByPasswordResetToken(string $token): ?User
    {
        if (!static::isPasswordResetTokenValid($token)) {
            return null;
        }

        return static::findOne([
            'password_reset_token' => $token,
            'status' => self::STATUS_ACTIVE
        ]);
    }

    /**
     * Finds user by verification email token
     * @param string $token verify email token
     * @return User|null
     */
    public static function findByVerificationToken(string $token): ?User
    {
        return static::findOne([
            'verification_token' => $token,
            'status' => self::STATUS_INACTIVE
        ]);
    }

    /**
     * Finds out if password reset token is valid
     * @param string $token password reset token
     * @return bool
     */
    public static function isPasswordResetTokenValid(string $token): bool
    {
        if (empty($token)) {
            return false;
        }

        $timestamp = (int)substr($token, strrpos($token, '_') + 1);
        $expire = Yii::$app->params['user.passwordResetTokenExpire'];
        return $timestamp + $expire >= time();
    }

    /**
     * @return int
     */
    public function getId(): int
    {
        return $this->getPrimaryKey();
    }

    /**
     * @return string
     */
    public function getAuthKey(): string
    {
        return $this->auth_key;
    }

    /**
     * @param string $authKey
     * @return bool
     */
    public function validateAuthKey($authKey): bool
    {
        return $this->getAuthKey() === $authKey;
    }

    /**
     * Validates password
     * @param string $password password to validate
     * @return bool if password provided is valid for current user
     */
    public function validatePassword(string $password): bool
    {
        return Yii::$app->security->validatePassword($password, $this->password_hash);
    }

    /**
     * Generates password hash from password and sets it to the model
     * @throws \yii\base\Exception
     */
    public function setPassword(string $password): void
    {
        $this->password_hash = Yii::$app->security->generatePasswordHash($password);
    }

    /**
     * Generates "remember me" authentication key
     * @throws \yii\base\Exception
     */
    public function generateAuthKey(): void
    {
        $this->auth_key = Yii::$app->security->generateRandomString();
    }

    /**
     * Generates new password reset token
     * @throws \yii\base\Exception
     */
    public function generatePasswordResetToken(): void
    {
        $this->password_reset_token = Yii::$app->security->generateRandomString() . '_' . time();
    }

    /**
     * @throws \yii\base\Exception
     */
    public function generateEmailVerificationToken(): void
    {
        $this->verification_token = Yii::$app->security->generateRandomString() . '_' . time();
    }

    /**
     * Removes password reset token
     */
    public function removePasswordResetToken(): void
    {
        $this->password_reset_token = null;
    }
}
