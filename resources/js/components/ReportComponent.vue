<template>
    <div class="container mt-5 mb-5">
        <h1>{{ vacancyName }}</h1>
        <div class="row mt-3">
            <div class="col-md-6">
                <h3>Количество вакансий</h3>
                <p>{{ vacanciesCount }}</p>
            </div>
            <div class="col-md-6">
                <h3>З/п</h3>
                Ср: {{ averageSalary }}<br>
                Мин: {{ minSalary }}<br>
                Макс: {{ maxSalary }}
            </div>
        </div>

        <div class="mt-5">
            <h3>Наиболее частые скиллы</h3>
            <ul class="list-group">
                <li v-for="(count, skill) in popularSkills"
                    class="list-group-item d-flex justify-content-between align-items-center">
                    {{ skill }}
                    <span class="badge badge-primary badge-pill">
                        {{ count }}
                    </span>
                </li>
            </ul>
        </div>

        <div class="mt-5">
            <h3>Список вакансий</h3>
            <div v-for="vacancy in vacancies" class="card">
                <div class="card-body">
                    <h5 class="card-title">{{ vacancy.vacancyName }}</h5>
                    <h6 class="card-subtitle mb-2 text-muted">{{ vacancy.companyName }}</h6>
                    <p class="card-text">{{ vacancy.salaryText }}</p>
                    <p v-for="skill in vacancy.skills" class="btn btn-primary mr-2">{{ skill }}</p>
                    <p>
                        <a :href="vacancy.link" class="card-link" target="_blank">Вакансия</a>
                    </p>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        data() {
            return {
                /** @var {String} vacancyName Вакансия */
                vacancyName: 'Php-программист',

                /** @var {Integer} vacanciesCount */
                vacanciesCount: 0,

                /** @var {Float} minSalary */
                minSalary: 0,

                /** @var {Float} maxSalary */
                maxSalary: 0,

                /** @var {Float} averageSalary */
                averageSalary: 0,

                /** @var {Mixed[]} vacancies Array of vacancies */
                vacancies: [],

                /** @var {String[]} popularSkills Array of the most often required skills in vacancies */
                popularSkills: []
            };
        },

        mounted() {
            axios.get('/api/v3/parser').then((response) => {
                for (let index in response.data.vacancies) {
                    this.vacancies.push(JSON.parse(response.data.vacancies[index]));
                }

                this.vacanciesCount = response.data.vacanciesCount;
                this.popularSkills = response.data.popularSkills;
                this.minSalary = response.data.salaries.minSalary;
                this.maxSalary = response.data.salaries.maxSalary;
                this.averageSalary = response.data.salaries.averageSalary;
            });
        }
    }
</script>
