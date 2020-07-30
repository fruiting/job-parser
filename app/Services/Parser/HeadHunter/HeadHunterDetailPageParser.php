<?php

namespace App\Services\Parser\HeadHunter;

use App\Services\Parser\ParserDetailBaseAbstract;

/**
 * Class HeadHunterDetailPageParser describes parser logic for hh.ru vacancy detail page
 *
 * @package App\Services\Parser\HeadHunter
 */
class HeadHunterDetailPageParser extends ParserDetailBaseAbstract
{
    /**
     * Parses specific vacancy info
     *
     * @return void
     *
     * @throws \PHPHtmlParser\Exceptions\ChildNotFoundException
     * @throws \PHPHtmlParser\Exceptions\NotLoadedException
     */
    public function loadVacancyInfo(): void
    {
        dd($this->dom->find('span'));
    }
}
