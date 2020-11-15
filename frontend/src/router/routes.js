import callbacks from './callbacks';


const routes = {
    home: {
        pathname: '/',
        callback: callbacks.renderHomePage,
    },
    login: {
        pathname: '/login',
        callback: callbacks.renderLoginPage,
    }
    dashboard: {
        pathname: '/dashboard',
        callback: callbacks.renderDashboardPage,
    },
}
