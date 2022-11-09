import hasRole from './permission/hasRole'
import hasPerms from './permission/hasPerms'
import copyText from './common/copyText'

export default function directive(app) {
  app.directive('hasRole', hasRole)
  app.directive('hasPerms', hasPerms)
  app.directive('copyText', copyText)
}
