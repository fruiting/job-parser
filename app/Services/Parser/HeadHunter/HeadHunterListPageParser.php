<?php

namespace App\Services\Parser\HeadHunter;

use App\Services\Parser\ParserListBaseAbstract;
use PHPHtmlParser\Dom\Node\Collection;
use PHPHtmlParser\Dom\Node\HtmlNode;
use PHPHtmlParser\Exceptions\ChildNotFoundException;
use PHPHtmlParser\Exceptions\NotLoadedException;
use Throwable;

/**
 * Class HeadHunterListPageParser describes parser logic for hh.ru vacancies list page
 *
 * @package App\Services\Parser\HeadHunter
 */
class HeadHunterListPageParser extends ParserListBaseAbstract
{
    /** @var string Link to parse */
    public const LINK = 'https://hh.ru/search/vacancy?area=1&st=searchVacancy&fromSearch=true&text=';

    /**
     * Parses count of vacancies
     *
     * @return void
     */
    public function loadVacanciesCount(): void
    {
        try {
            /** @var HtmlNode $html */
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
        $blocks = $this->dom->find('div.vacancy-serp-item');
        foreach ($blocks as $block) {
            try {
                /** @var Collection|HtmlNode[] $collection */
                $collection = $block->find('a');
                if ($collection[0] && $collection[0]->getAttribute('data-qa') == 'vacancy-serp__vacancy-title') {
                    $this->vacanciesUrls[] = $collection[0]->getAttribute('href');
                }
            } catch (Throwable $exception) {
                //todo log it!
            }
        }

        $this->vacanciesUrls = array_unique($this->vacanciesUrls);
    }
}
