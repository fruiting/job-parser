<?php

namespace App\Services\Parser;

/**
 * Interface GeneralParserInterface describes methods of general parser
 *
 * @package App\Services\Parser
 */
interface GeneralParserInterface
{
    /**
     * Returns object to parse list page
     *
     * @return ListPageParserInterface
     */
    public function getListPageParser(): ListPageParserInterface;

    /**
     * Returns object to parse detail page
     *
     * @return DetailPageParserInterface
     */
    public function getDetailPageParser(): DetailPageParserInterface;
}
