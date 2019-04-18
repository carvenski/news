import Taro, { Component } from '@tarojs/taro'
import { View, Text } from '@tarojs/components'
import './index.css'

export default class Index extends Component {

  config = {
    navigationBarTitleText: '首页'
  }

  componentWillMount () { }

  componentDidMount () { 
    console.log("==========")

    Taro.getSetting({
      complete: function (res) {
        console.log("=> 用户是否登陆", res.authSetting["scope.userInfo"])
      }
    });

    Taro.getUserInfo({
      scope: "scope.userInfo",
      complete: function (res) {
        console.log(res);
      }

    });


  }


  componentWillUnmount () { }

  componentDidShow () { }

  componentDidHide () { }

  render () {
    return (
      <View>
        <Text>Please Login First.</Text>
        <button open-type="getUserInfo">请您授权</button>
      </View>
    )
  }
}
