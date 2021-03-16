<?php
declare(strict_types=1);

namespace backend\assets;

use yii\web\AssetBundle;

/**
 * Class AppAsset
 * @package backend\assets
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
     * @var string[]
     */
    public $css = [
        'css/site.css',
    ];

    /**
     * @var array
     */
    public $js = [
    ];

    /**
     * @var string[]
     */
    public $depends = [
        'yii\web\YiiAsset',
        'yii\bootstrap\BootstrapAsset',
    ];
}
