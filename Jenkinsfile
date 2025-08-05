pipeline {
    agent any

    stages {
        stage('Deploy') {
            steps {
                script {
                    withCredentials([
                        string(credentialsId: 'PROJECT_DIR', variable: 'PROJECT_DIR'),
                        usernamePassword(credentialsId: 'GITHUB_TOKEN', usernameVariable: 'GIT_USER', passwordVariable: 'GIT_PASS')
                    ]) {
                        sh """
                            echo "Usando ruta: /srv/$PROJECT_DIR/appdata/libreria-mariela-api"
                            cd /srv/$PROJECT_DIR/appdata/libreria-mariela-api
                            git pull https://$GIT_USER:$GIT_PASS@github.com/brunotarditi/libreria-mariela-api.git main

                            cd /srv/$PROJECT_DIR/appdata/libreria_mariela_docker_compose
                            docker compose -f libreria_mariela_docker_compose.yml --env-file libreria_mariela_docker_compose.env down
                            docker compose -f libreria_mariela_docker_compose.yml --env-file libreria_mariela_docker_compose.env pull
                            docker compose -f libreria_mariela_docker_compose.yml --env-file libreria_mariela_docker_compose.env up -d --build
                        """
                    }
                }
            }
        }
    }
}
