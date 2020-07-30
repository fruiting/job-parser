<?php

namespace App\Services\Parser;

use PHPHtmlParser\Dom\Node\Collection;
use PHPHtmlParser\Dom\Node\HtmlNode;
use PHPHtmlParser\Exceptions\ChildNotFoundException;
use PHPHtmlParser\Exceptions\NotLoadedException;
use Throwable;

/**
 * Class HeadHunterParser describes logic of parsing hh.ru
 *
 * @package App\Services\Parser
 */
class HeadHunterParser extends ParserBaseAbstract
{
    /** @var string Web-site link to parse */
    public const LINK = 'https://hh.ru/search/vacancy?st=searchVacancy&text=PHP+%D0%BF%D1%80%D0%BE%D0%B3%D1%80%D0%B0%D0%BC%D0%BC%D0%B8%D1%81%D1%82&search_field=name&area=1&salary=&currency_code=RUR&experience=doesNotMatter&order_by=relevance&search_period=&items_on_page=50&no_magic=true&L_save_area=true&from=suggest_post';

    /**
     * Parses count of vacancies
     *
     * @return void
     */
    public function loadVacanciesCount(): void
    {
        /** @var HtmlNode $html */
        try {
            $html = $this->dom->find('h1');
            $header = $html->getChildren()[0];
            preg_match('!\d+!', $header->text(), $matches);
            $this->vacanciesCount = (int) $matches[0];
        } catch (ChildNotFoundException | NotLoadedException | Throwable $exception) {
            $this->vacanciesCount = 0;

            //todo log it!
        }
    }

    /**
     * Parses all vacancies for description
     *
     * @return void
     *
     * @throws ChildNotFoundException
     * @throws NotLoadedException
     */
    public function loadVacanciesInfo(): void
    {
        /** @var Collection|HtmlNode[] $blocks */
        $blocks = $this->dom->find('div');
        foreach ($blocks as $block) {
            try {
                if (preg_match('/vacancy-serp-item/', $block->getAttribute('class'))) {
                    /** @var Collection|HtmlNode[] $collection */
                    $collection = $block->find('a');
                    if ($collection[0]
                        && $collection[0]->getAttribute('data-qa') == 'vacancy-serp__vacancy-title') {
                        $this->vacanciesUrls[] = $collection[0]->getAttribute('href');
                    }
                }
            } catch (Throwable $exception) {
                //todo log it!
            }
        }

        $this->vacanciesUrls = array_unique($this->vacanciesUrls);
    }

    /**
     * Parses specific vacancy info
     *
     * @return void
     *
     * @throws ChildNotFoundException
     * @throws \PHPHtmlParser\Exceptions\CircularException
     * @throws \PHPHtmlParser\Exceptions\ContentLengthException
     * @throws \PHPHtmlParser\Exceptions\LogicalException
     * @throws \PHPHtmlParser\Exceptions\StrictException
     * @throws \Psr\Http\Client\ClientExceptionInterface
     */
    public function loadSpecificVacancyInfo(): void
    {
        foreach ($this->vacanciesUrls as $url) {
            $dom = DomHelper::getInitedDom($url);
            $dom->find('span');

            //todo похоже придется декомпозировать на классы для списка и деталки
        }
    }
}
