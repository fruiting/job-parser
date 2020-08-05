<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

/**
 * Class User describe model of user
 *
 * @package App\Models
 *
 * @property-read string $email
 */
class User extends Model
{
    /** @var string $table Table code */
    protected $table = 'users';

    /** @var string[] $fillable Fillable columns */
    protected $fillable = ['email'];

    /** @var bool $timestamps Use date_create and date_update fields */
    public $timestamps = false;
}
