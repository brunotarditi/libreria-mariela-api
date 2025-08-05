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
                        sh '''
                            echo "==> Desplegando en ruta: /srv/$PROJECT_DIR"

                            cd /srv/$PROJECT_DIR/appdata/libreria-mariela-api
                            echo "==> Haciendo pull del código..."
                            git pull https://$GIT_USER:$GIT_PASS@github.com/brunotarditi/libreria-mariela-api.git main

                            cd /srv/$PROJECT_DIR/appdata/libreria_mariela_docker_compose
                            echo "==> Apagando contenedores actuales..."
                            docker-compose -f libreria_mariela_docker_compose.yml --env-file libreria_mariela_docker_compose.env down

                            echo "==> Bajando imágenes actualizadas (si corresponde)..."
                            docker-compose -f libreria_mariela_docker_compose.yml --env-file libreria_mariela_docker_compose.env pull

                            echo "==> Reconstruyendo imágenes y levantando contenedores..."
                            docker-compose -f libreria_mariela_docker_compose.yml --env-file libreria_mariela_docker_compose.env up -d --build
                        '''
                    }
                }
            }
        }
    }
}
