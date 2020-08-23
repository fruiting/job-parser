<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

/**
 * Class Vacancy describes model of vacancy
 *
 * @package App\Models
 *
 * @property-read int       $id
 * @property-read string    $name
 */
class Vacancy extends Model
{
    /** @var string $table Table code */
    protected $table = 'vacancies';

    /** @var string[] $fillable Fillable columns */
    protected $fillable = ['name'];

    /** @var bool $timestamps Use date_create and date_update fields */
    public $timestamps = false;
}
