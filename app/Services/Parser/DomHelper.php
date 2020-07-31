<?php

namespace App\Services\Parser;

use PHPHtmlParser\Dom;

/**
 * Class DomHelper describes initiation of Dom object
 *
 * @package App\Services\Parser
 */
class DomHelper
{
    /**
     * Returns inited dom object by link
     *
     * @param string $link
     *
     * @return Dom
     *
     * @throws \PHPHtmlParser\Exceptions\ChildNotFoundException
     * @throws \PHPHtmlParser\Exceptions\CircularException
     * @throws \PHPHtmlParser\Exceptions\ContentLengthException
     * @throws \PHPHtmlParser\Exceptions\LogicalException
     * @throws \PHPHtmlParser\Exceptions\StrictException
     * @throws \Psr\Http\Client\ClientExceptionInterface
     */
    public static function getInitedDom(string $link): Dom
    {
        $dom = new Dom();
        $dom->loadFromUrl($link);

        return $dom;
    }
}
