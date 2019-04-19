import Taro, { Component } from '@tarojs/taro'
import { View, Text } from '@tarojs/components'
import './index.css'

export default class Index extends Component {

  config = {
    navigationBarTitleText: '首页'
  }

  componentWillMount () { }

  // 页面加载后开始执行js逻辑
  componentDidMount () { 
    // 判断是否之前登录过
    Taro.getSetting({
      complete: function (res) {
        is_login = res.authSetting["scope.userInfo"]
        console.log("is_login = ", is_login)
        if (is_login) {
          // 已登录,获取user info
          Taro.getUserInfo({
            scope: "scope.userInfo",
            complete: function (res) {
              console.log("=> got user info:");
              console.log(res);
            }
          });
        } else {
          // 未登录,提示授权
          Taro.showToast({
            title: '请您先授权',
            icon: 'success',
            duration: 3000,
            mask: true
          }).then(res => console.log(res))
        }
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
        <button open-type="getUserInfo">授权按钮</button>
      </View>
    )
  }
}
