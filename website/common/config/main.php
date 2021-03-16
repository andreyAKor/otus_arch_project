<?php
declare(strict_types=1);

return [
    'id' => 'otus_arch_project',
    'name' => 'OTUS Arch Project',
    'charset' => 'utf-8', // Кодировка системы
    'sourceLanguage' => 'ru', // Локализация
    'language' => 'ru', // Язык системы
    'timeZone' => 'Europe/Moscow',
    'aliases' => [
        '@bower' => '@vendor/bower-asset',
        '@npm' => '@vendor/npm-asset'
    ],
    'vendorPath' => dirname(dirname(__DIR__)) . '/vendor',

    // Точка входа по умолчанию
    'defaultRoute' => 'site/index',

    'components' => [
        'cache' => [
            'class' => \yii\caching\FileCache::class
        ],

        // Устанавливаем таймзону
        'formatter' => [
            'class' => \yii\i18n\Formatter::class,
            'timeZone' => 'Europe/Moscow',
            'locale' => 'ru-RU'
        ],

        // Логгер monolog
        'monolog' => [
            'class' => \common\components\monolog\MonologComponent::class,
            'channels' => [
                'main' => [
                    'handler' => [
                        [
                            'type' => 'rotating_file',
                            'path' => '@app/runtime/logs/log_' . date('Y-m-d') . '.log'
                        ]
                    ]
                ],
                'services' => [
                    'handler' => [
                        [
                            'class' => \common\components\monolog\handler\Api::class,
                            'table' => '{{%services_api_Log}}',
                            'logLiferime' => 24 * 3
                        ]
                    ]
                ]
            ]
        ],

        // Компонент для работы с api сервисов
        'services' => [
            'class' => \common\components\services\Service::class
        ]
    ]
];
