var appUserList = new Vue({
    el: '#app-user-list',
    data: {
        userList: [],
        updateName: '',
        updateAge: '',
        updateWeight: '',
        name: '',
        age: '',
        weight: '',
    },
    mounted: function() {
        this.queryAllUser();
    },
    methods: {
        queryAllUser: function() {
            let self = this

            this.userList = []

            fetch('/query/user/')
                .then(function(response) {
                    return response.json();
                })
                .then(function(myJson) {
                    if (myJson["error-code"] !== 0) {
                        alert("error-code:" + myJson["error-code"] + ", " + "error-text:" + myJson["error-text"]);
                        return;
                    }
                    self.userList = myJson.data
                });
        },
        cancelAdd: function() {
            this.name = ''
            this.age = ''
            this.weight = ''
        },
        cancelUpdate: function() {
            this.updateName = ''
            this.updateAge = ''
            this.updateWeight = ''
        },
        updateUser: function() {
            let self = this

            let params = "?name=" + self.updateName
            let url = "/update/user/" + params

            fetch(url, {
                    body: JSON.stringify({
                        name: self.updateName,
                        age: parseInt(self.updateAge, 10),
                        weight: parseInt(self.updateWeight, 10),
                    }),
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    method: 'POST',
                })
                .then(function(response) {
                    return response.json();
                })
                .then(function(myJson) {
                    if (myJson["error-code"] !== 0) {
                        alert("error-code:" + myJson["error-code"] + ", " + "error-text:" + myJson["error-text"]);
                    }
                    appUserList.$emit('updateUserList')
                });
        },
        delUser: function(name) {
            appUserList.$emit('delUser', name)
        },
        setUpdateDivInfo: function(name, age, weight) {
            appUserList.$emit('setUpdateDiv', name, age, weight)
        },
        addUser: function() {
            let self = this

            fetch("/create/user/", {
                    body: JSON.stringify({
                        name: self.name,
                        age: parseInt(self.age, 10),
                        weight: parseInt(self.weight, 10),
                    }),
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    method: 'POST',
                })
                .then(function(response) {
                    return response.json();
                })
                .then(function(myJson) {
                    if (myJson["error-code"] !== 0) {
                        alert("error-code:" + myJson["error-code"] + ", " + "error-text:" + myJson["error-text"]);
                    }
                    appUserList.$emit('updateUserList')
                });
        },
    },
})

appUserList.$on('updateUserList', function() {
    appUserList.queryAllUser();
})

appUserList.$on('setUpdateDiv', function(name, age, weight) {
    appUserList.updateName = name
    appUserList.updateAge = age
    appUserList.updateWeight = weight
})

appUserList.$on('delUser', function(name) {
    let self = this

    let params = "?name=" + name
    let url = "/del/user/" + params

    console.log("del user, url:" + url)

    fetch(url, {
            headers: {
                'Content-Type': 'application/json'
            },
            method: 'POST',
        })
        .then(function(response) {
            return response.json();
        })
        .then(function(myJson) {
            if (myJson["error-code"] !== 0) {
                alert("error-code:" + myJson["error-code"] + ", " + "error-text:" + myJson["error-text"]);
            }
            appUserList.$emit('updateUserList')
        });
})