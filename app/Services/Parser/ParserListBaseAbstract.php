<?php

namespace App\Services\Parser;

use App\Services\Parser\DomHelper;
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
     * @return void
     */
    public function execute(): void
    {
        $this->dom = DomHelper::getInitedDom(static::LINK);
        $this->loadVacanciesCount();
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
