pipeline {
    agent any

    parameters {
        string(name: 'PROJECT_DIR', defaultValue: '', description: 'Ruta protegida del proyecto')
    }

    stages {
        stage('Deploy') {
            steps {
                script {
                    sh '''
                        cd $PROJECT_DIR/libreria-mariela-api
                        git pull origin main
                        cd $PROJECT_DIR/libreria_mariela_docker_compose
                        docker compose -f libreria_mariela_docker_compose.yml --env-file libreria_mariela_docker_compose.env down
                        docker compose -f libreria_mariela_docker_compose.yml --env-file libreria_mariela_docker_compose.env pull
                        docker compose -f libreria_mariela_docker_compose.yml --env-file libreria_mariela_docker_compose.env up -d --build
                    '''
                }
            }
        }
    }
}
