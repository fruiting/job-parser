<?php

namespace App\Services\Parser;

use PHPHtmlParser\Dom;

/**
 * Class ParserDetailBaseAbstract describes parse logic of detail page
 *
 * @package App\Services\Parser
 */
abstract class ParserDetailBaseAbstract implements DetailPageParserInterface
{
    /** @var Dom $dom Dom parser object */
    protected $dom;

    /**
     * Parses vacancy detail page
     *
     * @param string $link Vacancy page link
     *
     * @return void
     *
     * @throws \PHPHtmlParser\Exceptions\ChildNotFoundException
     * @throws \PHPHtmlParser\Exceptions\CircularException
     * @throws \PHPHtmlParser\Exceptions\ContentLengthException
     * @throws \PHPHtmlParser\Exceptions\LogicalException
     * @throws \PHPHtmlParser\Exceptions\NotLoadedException
     * @throws \PHPHtmlParser\Exceptions\StrictException
     * @throws \Psr\Http\Client\ClientExceptionInterface
     */
    public function execute(string $link): void
    {
        $this->dom = DomHelper::getInitedDom($link);
        $this->loadVacancyInfo();
    }
}
