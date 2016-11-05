
  Vue.component('my-login', {
      template: '#login',
      data() {
        return {
          name: 'alex',
          password: '1234'
        }
      },
      created:()=>{
        console.log('my-login')
      },
      methods: {
        onLogin() {
          var self=this
         // console.log(self.password)
       
          axios.post('/login', { name:self.name ,password:self.password})
            .then((data ) => {
                console.log(data)
             })
            self.$emit('pass')
          }  
      }
    })
Vue.component('add-user', {
      template: '#adduser',
      data() {
        return {
          
          name: 'alex',
        }
      },
      props: {
         users: {
           type: Array,
          default() {
            return []
          }
        }
      },
      created:()=>{
        console.log('add-user')
      },
      methods: {
         onAddUser() {
           let self=this
           axios.post('/api/users', { name:self.name }).then(({data} ) => {
             self.users.push(data.name)
            // this.exclamations = [data.exclamation].concat(this.exclamations);
             console.log(data)
           })
          this.name = ''
        }
      }
    })

var demo = new Vue({
  el: '#app',
  data() {
      return {
        name: '',
        password: '',
        showLogin:true,
        users:[],
      }
  },
  created:()=>{
    console.log('demo')
  },
  methods: {
    hideLogin () {
      this.showLogin=false
       console.log('hide')
    },
   
      
  }
})
