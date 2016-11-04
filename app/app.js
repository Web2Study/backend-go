var apiURL = 'http://localhost:3000/users'

/**
 * Actual demo
 */

var demo = new Vue({

  el: '#demo',

  data: {
 
    msg: "hello"
  },

  created: function () {
    //this.fetchData()
  },


  methods: {
    fetchData: function () {
      var xhr = new XMLHttpRequest()
      var self = this
      xhr.open('GET', apiURL + self.currentBranch)
      xhr.onload = function () {
       // self.commits = JSON.parse(xhr.responseText)
        console.log(self.commits[0].html_url)
      }
      xhr.send()
    }
  }
})
