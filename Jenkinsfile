pipeline {
    agent any

    stages {
        stage('Deploy') {
            steps {
                script {
                    withCredentials([string(credentialsId: 'PROJECT_DIR', variable: 'PROJECT_DIR')]) {
                        sh '''
                            echo "Usando ruta: $PROJECT_DIR"
                            cd /srv/$PROJECT_DIR/appdata/libreria-mariela-api
                            git pull origin main
                            cd /srv/$PROJECT_DIR/appdata/libreria_mariela_docker_compose
                            docker compose -f libreria_mariela_docker_compose.yml --env-file libreria_mariela_docker_compose.env down
                            docker compose -f libreria_mariela_docker_compose.yml --env-file libreria_mariela_docker_compose.env pull
                            docker compose -f libreria_mariela_docker_compose.yml --env-file libreria_mariela_docker_compose.env up -d --build
                        '''
                    }
                }
            }
        }
    }
}
