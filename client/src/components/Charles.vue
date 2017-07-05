<template>
    <div class="columns justify-center">
        <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop">
            <div class="column">
                <div class="control has-icon has-icon-right">
                    <span class="icon is-small is-right is-clickable-icon" @click="onClickIcon">
                        <i class="fa fa-times-circle"></i>
                    </span>
                    <input v-model="input" class="input" placeholder="请输入用户名">
                </div>
            </div>
            <div class="column">
                <input v-model="key" class="input" placeholder="注册码在这里" readonly>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        name: 'Charles',
        data() {
            return {
                input: '',
                key: '',
            }
        },
        methods: {
            onClearName() {
                this.input = ''
            },
            onClickIcon() {
                this.input = ''
            },
        },
        subscriptions() {
            return {
                key: this.$watchAsObservable('input').
                    debounceTime(200).
                    distinctUntilChanged().
                    pluck('newValue').
                    switchMap(name => name ? fetchCharlesKey(name) : Rx.Observable.of({key: ''})).
                    map(data => data.key).
                    startWith(''),
            }
        },
        created() {
            this.$emit('created', 'Charles')
        },
    }

    /**
     *
     * @param { string }name
     * @return {*|Observable<T>}
     */
    function fetchCharlesKey(name) {
        let p = fetch('/api/keygen/charles', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({name: name}),
        }).then(response => response.json()).catch(() => 'Internal Error')
        return Rx.Observable.fromPromise(p)
    }
</script>
