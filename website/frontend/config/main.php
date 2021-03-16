<?php
declare(strict_types=1);

$params = array_merge(
    require __DIR__ . '/../../common/config/params.php',
    require __DIR__ . '/../../common/config/params-local.php',
    require __DIR__ . '/params.php',
    require __DIR__ . '/params-local.php'
);

return [
    'id' => 'otus_arch_project-frontend',
    'name' => 'OTUS Arch Project',
    'charset' => 'utf-8', // Кодировка системы
    'sourceLanguage' => 'ru-RU', // Локализация
    'language' => 'ru-RU', // Язык системы
    'timeZone' => 'Europe/Moscow',
    'basePath' => dirname(__DIR__),
    'bootstrap' => [
        'log'
    ],
    'controllerNamespace' => 'frontend\controllers',
    'components' => [
        'request' => [
            'csrfParam' => '_csrf-frontend',
            'baseUrl' => ''
        ],
        'user' => [
            'identityClass' => 'common\models\User',
            'enableAutoLogin' => true,
            'identityCookie' => [
                'name' => '_identity-frontend',
                'httpOnly' => true
            ]
        ],
        'session' => [
            // this is the name of the session cookie used for login on the frontend
            'name' => 'advanced-frontend'
        ],
        'log' => [
            'traceLevel' => YII_DEBUG ? 3 : 0,
            'targets' => [
                [
                    'class' => \yii\log\FileTarget::class,
                    'levels' => [
                        'error',
                        'warning'
                    ]
                ]
            ]
        ],
        'errorHandler' => [
            'errorAction' => 'site/error'
        ],
        'urlManager' => [
            'enablePrettyUrl' => true,
            'showScriptName' => false,
            'rules' => [
                '' => 'site/index',
                [
                    'pattern' => '/bid/<bidId:[\d\w\-]+>',
                    'route' => 'site/view',
                    'suffix' => '.html'
                ],
            ],
        ],
        'i18n' => [
            'translations' => [
                'common*' => [
                    'class' => \yii\i18n\PhpMessageSource::class,
                    'basePath' => '@app/messages'
                ]
            ]
        ],
        'assetManager' => [
            'basePath' => '@webroot/assets',
            'baseUrl' => '@web/assets'
        ]
    ],
    'params' => $params
];
