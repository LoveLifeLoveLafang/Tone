import axios from 'axios'

axios.get('../../static/swiper.json').then(function(response){
    console.log(response)
}).catch(function(err){
    console.log(err)
})
