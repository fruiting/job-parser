<?php

namespace App\Services\Parser\HeadHunter;

use App\Services\Parser\ParserDetailBaseAbstract;
use phpDocumentor\Reflection\Types\Collection;
use PHPHtmlParser\Dom\Node\HtmlNode;

/**
 * Class HeadHunterDetailPageParser describes parser logic for hh.ru vacancy detail page
 *
 * @package App\Services\Parser\HeadHunter
 */
class HeadHunterDetailPageParser extends ParserDetailBaseAbstract
{
    /**
     * Loads vacancy name
     *
     * @return void
     *
     * @throws \PHPHtmlParser\Exceptions\ChildNotFoundException
     * @throws \PHPHtmlParser\Exceptions\NotLoadedException
     */
    public function loadVacancyName(): void
    {
        $header = $this->dom->find('h1[data-qa="vacancy-title"]');
        $this->vacancyName = $header[0]->getChildren()[0]->text();
    }

    /**
     * Loads salary info
     *
     * @return void
     *
     * @throws \PHPHtmlParser\Exceptions\ChildNotFoundException
     * @throws \PHPHtmlParser\Exceptions\CircularException
     * @throws \PHPHtmlParser\Exceptions\NotLoadedException
     */
    public function loadSalary(): void
    {
        /** @var Collection|HtmlNode[] $paragraphs */
        $paragraphs = $this->dom->find('p.vacancy-salary');
        foreach ($paragraphs as $paragraph) {
            $this->salaryText = $paragraph->getChildren()[0]->getChildren()[0]->text();
            preg_match_all('/(\+\d+)?\s*(\(\d+\))?([\s-]?\d+)+/', $this->salaryText, $matches);
            $this->salaryRange = array_map(function (string $salary) {
                return (float) str_replace(' ', '', $salary);
            }, $matches[0]);
            break;
        }
    }

    /**
     * Loads company name
     *
     * @return void
     *
     * @throws \PHPHtmlParser\Exceptions\ChildNotFoundException
     * @throws \PHPHtmlParser\Exceptions\CircularException
     * @throws \PHPHtmlParser\Exceptions\NotLoadedException
     */
    public function loadCompany(): void
    {
        /** @var Collection|HtmlNode[] $links */
        $links = $this->dom->find('a[data-qa="vacancy-company-name"]');
        foreach ($links as $link) {
            $this->companyName = $link->getChildren()[0]->getChildren()[0]->text();
            break;
        }
    }

    /**
     * Loads requirement skills of vacancy
     *
     * @return void
     *
     * @throws \PHPHtmlParser\Exceptions\ChildNotFoundException
     * @throws \PHPHtmlParser\Exceptions\NotLoadedException
     */
    public function loadSkills(): void
    {
        $this->skills = [];

        /** @var Collection|HtmlNode[] $divs */
        $divs = $this->dom->find('span.bloko-tag__section');
        foreach ($divs as $div) {
            $this->skills[] = $div->getChildren()[0]->text();
        }
    }
}
