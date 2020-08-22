<template>
    <div class="card">
        <div class="card-body">
            <h5 class="card-title">{{ vacancy.vacancyName }}</h5>
            <h6 class="card-subtitle mb-2 text-muted">{{ vacancy.companyName }}</h6>
            <p class="card-text">{{ vacancy.salaryText }}</p>
            <p v-for="skill in vacancy.skills" class="btn btn-primary mr-2">
                {{ skill }}
            </p>
            <p>
                <a :href="vacancy.link" class="card-link" target="_blank">Вакансия</a>
            </p>
        </div>
    </div>
</template>

<script>
    export default {
        props: ['vacancy']
    }
</script>
