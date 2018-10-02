<template>
  <div class="recommend">
    <div class="recommend-content">
      <div class="slider-wrapper" v-if="recommends.length">
        <slider>
          <div v-for="(item, index) in recommends" :key="item.id">
            <a :href="item.linkUrl">
              <img :src="picList[index]" :alt="picList[index]">
            </a>
          </div>
        </slider>
      </div>
      <div class="recommend-list">
        <h1 class="list-title">热门歌单推荐</h1>
        <ul>
            <li>
                <div>
                    <img src="" alt="">
                </div>
                <div>
                    <h2></h2>
                    <p></p>
                </div>
            </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'
import Slider from 'base/slider/slider'

Vue.use(VueAxios, axios)

export default {
  data () {
    return {
      recommends: [],
      picList: []
    }
  },
  created () {
    this._getRecommends()
    this._getPicSlider()
    this._getSonglist()
  },
  methods: {
    _getRecommends () {
      this.axios.get('static/swiper.json').then(response => {
        this.recommends = response.data.data.slider
      })
    },
    _getPicSlider () {
      this.axios.get('static/sliderInfo.json').then(response => {
        let picList = response.data
        for (let i = 0; i < picList.length; i++) {
            picList[i] = 'static/' + picList[i]
        }
        this.picList = picList
      })
    },
    _getSonglist () {
        
    }
  },
  components: {
    Slider
  }
}
</script>

<style>
div {
  color: #ce7692;
}

.recommend {
    position: fixed;
    width: 100%;
    top: 88px;
    bottom: 0;
}

.recommend-content {
    height: 100%;
    overflow: hidden;
}

.slider-wrapper {
    position: relative;
    width: 100%;
    overflow: hidden;
}
</style>
