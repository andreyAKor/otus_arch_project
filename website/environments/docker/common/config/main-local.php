<?php
declare(strict_types = 1);

return [
    'components' => [
        'db' => [
            'class' => 'yii\db\Connection',
            'dsn' => 'mysql:host=mysql;port=3306;dbname=otus_arch_project',
            'username' => 'root',
            'password' => 'Qybf0H8aQ0',
            'charset' => 'utf8'
        ],
        'mailer' => [
            'class' => 'yii\swiftmailer\Mailer',
            'viewPath' => '@common/mail',
            // send all mails to a file by default. You have to set
            // 'useFileTransport' to false and configure a transport
            // for the mailer to send real emails.
            'useFileTransport' => true
        ],

        // Компонент для работы с api сервисов
        'services' => [
            'url' => 'http://0.0.0.0:6080'
        ]
    ]
];
