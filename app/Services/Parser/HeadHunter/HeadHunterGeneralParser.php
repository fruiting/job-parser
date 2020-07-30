<?php

namespace App\Services\Parser\HeadHunter;

use App\Services\Parser\DetailPageParserInterface;
use App\Services\Parser\GeneralParserInterface;
use App\Services\Parser\ListPageParserInterface;

/**
 * Class HeadHunterGeneralParser describes general
 *
 * @package App\Services\Parser\HeadHunter
 */
class HeadHunterGeneralParser implements GeneralParserInterface
{
    /**
     * Returns object to parse list page
     *
     * @return ListPageParserInterface
     */
    public function getListPageParser(): ListPageParserInterface
    {
        return new HeadHunterListPageParser();
    }

    /**
     * Returns object to parse detail page
     *
     * @return DetailPageParserInterface
     */
    public function getDetailPageParser(): DetailPageParserInterface
    {
        return new HeadHunterDetailPageParser();
    }
}
