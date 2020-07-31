<?php

namespace App\Services\Parser;

use PHPHtmlParser\Dom;

/**
 * Class ParserListBaseAbstract describes base logic of list page parser
 *
 * @package App\Services\Parser
 */
abstract class ParserListBaseAbstract implements ListPageParserInterface
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
     * @param string $vacancyTitle Title to search
     * @param int $page Page to parse
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
    public function execute(string $vacancyTitle, int $page): void
    {
        $this->dom = DomHelper::getInitedDom(static::LINK . $vacancyTitle . '&page=' . $page);
        $this->loadVacanciesInfo();
    }

    /**
     * Returns count of vacancies
     *
     * @return int
     */
    public function getVacanciesCount(): int
    {
        return $this->vacanciesCount;
    }

    /**
     * Returns array of urls of vacancies
     *
     * @return string[]
     */
    public function getVacanciesUrls(): array
    {
        return $this->vacanciesUrls;
    }
}
