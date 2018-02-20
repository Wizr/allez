<template>
    <div class="columns justify-center">
        <div class="column is-narrow" v-if="accounts">
            <div class="column is-narrow loading" v-if="accounts.error">唉呀，出错了☹️</div>
            <div v-else>
                <div class="box"  v-for="(account, index) in accounts" :key="index">
                    <span>{{index}}</span>
                    <span class="account-field">{{account['AppleID']}}</span>
                    <span class="account-field">{{account['Password']}}</span>
                </div>
            </div>
        </div>
        <span class="column is-narrow loading" v-else>正在获取，请稍等...</span>
    </div>
</template>

<script>
    import axios from 'axios'
    export default {
        name: 'AppStore',
        asyncComputed: {
            accounts() {
                return axios.get('/api/keygen/appstore').then(response => {
                    let err = undefined
                    let data = response.data
                    if (data) {
                        if (data['result']) {
                            return data['result']
                        }
                        err = data
                    } else {
                        err = {error: 'Server internal error!'}
                    }
                    console.log(err)
                    return err
                }).catch((e) => {
                    let err = {error: e.message}
                    console.log(err)
                    return err
                })
            },
        },
        created() {
            this.$emit('created', 'AppStore')
        },
    }
</script>

<style scoped>
    .loading {
        margin-top: 40px;
    }

    .box:not(:last-child) {
        margin-top: 20px;
    }

    .account-field {
        margin-left: 30px;
    }
</style>
