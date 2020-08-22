<?php

namespace App\Services\Parser;

/**
 * Class ParserGeneralBase describes general parser base logic
 *
 * @package App\Services\Parser
 */
abstract class ParserGeneralBase implements GeneralParserInterface
{
    /**
     * Returns pages count
     *
     * @param string $vacanciesCount
     *
     * @return int
     */
    public function getPagesCount(string $vacanciesCount): int
    {
        return (int) ceil($vacanciesCount / static::PAGE_SIZE);
    }
}
