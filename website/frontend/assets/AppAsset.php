<?php
declare(strict_types=1);

namespace frontend\assets;

use yii\web\AssetBundle;

/**
 * Class AppAsset
 * @package frontend\assets
 */
class AppAsset extends AssetBundle
{
    /**
     * @var string
     */
    public $basePath = '@webroot';

    /**
     * @var string
     */
    public $baseUrl = '@web';

    /**
     * @var array
     */
    public $css = [
        'https://unpkg.com/flickity@2/dist/flickity.min.css',
        'css/main.css'
    ];

    /**
     * @var array
     */
    public $js = [
        'https://unpkg.com/flickity@2/dist/flickity.pkgd.min.js',
        'scripts/app.js'
    ];

    /**
     * @var array
     */
    public $depends = [
        \yii\web\JqueryAsset::class,
        \yii\web\YiiAsset::class,
        \yii\bootstrap\BootstrapAsset::class
    ];
}
