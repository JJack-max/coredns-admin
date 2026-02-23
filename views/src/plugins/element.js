import Vue from 'vue'
import { Container, Header, Aside, Main, Menu, Submenu, MenuItem, Table, TableColumn, /* Form, FormItem, */ Input, Select, Option, Button, Alert, Message, Popover, Popconfirm } from 'element-ui'

Vue.use(Container)
Vue.use(Header)
Vue.use(Aside)
Vue.use(Main)
Vue.use(Menu)
Vue.use(Submenu)
Vue.use(MenuItem)
Vue.use(Table)
Vue.use(TableColumn)
// Vue.use(Form)
// Vue.use(FormItem)
Vue.use(Input)
Vue.use(Select)
Vue.use(Option)
Vue.use(Alert)
Vue.use(Button)
Vue.use(Popover)
Vue.use(Popconfirm)

Vue.prototype.$message = Message
