<?php

namespace App\Services\Parser;

use PHPHtmlParser\Dom;

/**
 * Class ParserBaseAbstract describes base logic of parser
 *
 * @package App\Services\Parser
 */
abstract class ParserBaseAbstract implements ParserInterface
{
    /** @var Dom $dom Dom parser object */
    protected $dom;

    /** @var int $vacanciesCount Count of vacancies by title */
    protected $vacanciesCount;

    /** @var array|string[] $vacanciesUrls Array of detail pages of vacancies */
    protected $vacanciesUrls = [];

    /**
     * Executes parser
     *
     * @return void
     *
     * @throws \PHPHtmlParser\Exceptions\ChildNotFoundException
     * @throws \PHPHtmlParser\Exceptions\CircularException
     * @throws \PHPHtmlParser\Exceptions\ContentLengthException
     * @throws \PHPHtmlParser\Exceptions\LogicalException
     * @throws \PHPHtmlParser\Exceptions\StrictException
     * @throws \Psr\Http\Client\ClientExceptionInterface
     */
    public function execute(): void
    {
        $this->dom = DomHelper::getInitedDom(static::LINK);

        $this->loadVacanciesCount();
        $this->loadVacanciesInfo();
        $this->loadSpecificVacancyInfo();
    }
}
