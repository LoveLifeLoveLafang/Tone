<template>
    <div class="slider" ref="slider">
        <div class="slider-group" ref="sliderGroup">
            <slot></slot>
        </div>
        <div class="dots">
            <span class="dot" v-for="(item, index) in dots" :class="{active: currentPageIndex === index}"></span>
        </div>
    </div>
</template>

<script type="text/ecmascript-6">
import BScroll from 'better-scroll'
import {addClass} from 'common/js/dom'

export default {
  name: "slider",
  props: {
    loop: {
      type: Boolean,
      default: true
    },
    autoPlay: {
      type: Boolean,
      default: true
    },
    interval: {
      type: Number,
      default: 4000
    }
  },
  data () {
    return {
      dots: [],
      currentPageIndex: 0
    }
  },
  mounted() {
    setTimeout(() => {
      this._setSliderWidth(false)
      this._initDots()
      this._initSlider()

      if (this.autoPlay) {
          this._play()
      }
    }, 20)

    window.addEventListener('resize', () => {
        if (!this.slider) {
            return
        }

        this._setSliderWidth(true)
        this.slider.refresh()
    })
  },
  methods: {
      _setSliderWidth(isResize) {
          this.children = this.$refs.sliderGroup.children
          let sliderWidth = this.$refs.slider.clientWidth
          let width = 0

          for (let i = 0; i < this.children.length; i++) {
              width += sliderWidth
              let child = this.children[i]
              addClass(child, 'slider-item')
              child.style.width = sliderWidth + 'px'
          }

          if (this.loop && !isResize) {
              width += 2 * sliderWidth
          }

          this.$refs.sliderGroup.style.width = width + 'px'
      },
      _initSlider() {
          this.slider = new BScroll(this.$refs.slider,{
              scrollX: true,
              scrollY: false,
              momentum: false,
              snap: {
                  loop: this.loop,
                  threshold: 0.3,
                  speed: 1000
              },
          })

          //当轮播图滑动时监听滑动结束事件,从初始化的时候就开始监听
          //每一张图片滑动结束的事件
          this.slider.on('scrollEnd', () => {
              let pageIndex = this.slider.getCurrentPage().pageX
              this.currentPageIndex = pageIndex

              if (this.autoPlay) {
                  this._play()
              }
          })
      },
      _initDots() {
          this.dots = new Array(this.children.length)
      },
      //_play()的周期是等待4秒到下1个图片开始等待的事件，中间包含滑动需要的时间
      _play() {
          //this.currentPageIndex 从0开始
          clearTimeout(this.timer)
          this.timer = setTimeout(() => {
              this.slider.next()
          }, this.interval)
      }
  },
  destroyed () {
    clearTimeout(this.timer)
  }
}
</script>

<style type="text/css">
.slider {
    min-height: 1px;
}

.slider-group {
    position: relative;
    padding: 0px;
    overflow: hidden;
}

.slider-group .slider-item{
    float: left;
    box-sizing: border-box;
    overflow: hidden;
    text-align: center;
}

.slider-item a {
    display: block;
    width: 100%;
    text-decoration: none;
}

.slider-item a img {
    display: block;
    width: 100%;
}

.slider .dots {
    position: absolute;
    bottom: 12px;
    left: 0;
    right: 0;
    text-align: center;
    font-size: 0;
}

.dot {
    display: inline-block;
    margin-right: 4px;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background-color: rgba(0,0,0,0.5);
}

.active {
    width: 20px;
    border-radius: 5px;
    background-color: rgba(255, 255, 255, 0.3);
}
</style>
