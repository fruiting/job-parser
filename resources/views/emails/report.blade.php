@php
    /** Report link mail template */
    /** @var string $link Report link */
    /** @var \App\Models\Vacancy $vacancy */
@endphp

Информация о вакансии {{ $vacancy->name }} собрана. Отчет доступен по <a href="{{ $link }}">ссылке</a>.
